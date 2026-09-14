#!/usr/bin/env python3
import json,os,pathlib,re,shutil,signal,subprocess,sys,time
root=pathlib.Path.cwd();module=sys.argv[1];cwd=root if module=='root' else root/'evmd';out=root/'versioned-evidence';out.mkdir(exist_ok=True);results={}
def run(label,args):
 minimum=shutil.disk_usage(root).free;start=time.monotonic();reason=None
 with (out/(label+'.log')).open('w') as log:
  p=subprocess.Popen(args,cwd=cwd,stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
  while p.poll() is None:
   minimum=min(minimum,shutil.disk_usage(root).free)
   mem=int(next(l.split()[1] for l in pathlib.Path('/proc/meminfo').read_text().splitlines() if l.startswith('MemAvailable:')))*1024
   if minimum<5.25*1024**3 or mem<768*1024**2 or time.monotonic()-start>5400:
    reason='resource guard or stage deadline';os.killpg(p.pid,signal.SIGTERM)
    try:p.wait(timeout=10)
    except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
    break
   time.sleep(.5)
 results[label]={'command':args,'exit_code':p.returncode,'blocked':reason,'minimum_free_bytes':minimum,'seconds':round(time.monotonic()-start,2)}
 (out/'results.json').write_text(json.dumps(results,indent=2));print(label,json.dumps(results[label]),flush=True)
 if p.returncode:print((out/(label+'.log')).read_text()[-5000:])
 return p.returncode==0 and reason is None
run('modules',['go','list','-m','-json','all'])
run('provenance',['python3',str(root/'scripts/verify-go-fork-provenance.py')])
run('verify-sums',['go','mod','verify'])
run('build',['go','build','./...'])
run('tests',['go','test','-tags=test','./...'])
run('production-graph',['go','list','-deps','./...'])
packages=set((out/'production-graph.log').read_text().splitlines())
forbidden=sorted(p for p in packages if p.startswith('golang.org/x/crypto/openpgp') or p.startswith('github.com/pion/dtls/v2') or p.startswith('github.com/cosmos/cosmos-sdk/x/crisis') or p.startswith('cosmossdk.io/x/crisis') or p.startswith('github.com/cosmos/cosmos-sdk/contrib/x/crisis'))
raw_key_files=[str(p.relative_to(root)) for p in (root/'rpc').rglob('*.go') if re.search(r'ImportRawKey|personal_importRawKey',p.read_text())]
(out/'production-graph-assessment.json').write_text(json.dumps({'forbidden_production_packages':forbidden,'forbidden_raw_key_rpc_files':raw_key_files},indent=2))
if run('install-scanner',['go','install','golang.org/x/vuln/cmd/govulncheck@v1.7.0']):
 scanner=str(pathlib.Path(subprocess.check_output(['go','env','GOPATH'],cwd=cwd,text=True).strip())/'bin/govulncheck')
 run('source-scan',[scanner,'-json','./...'])
 if module=='evmd':
  head=os.environ['QUALIFICATION_HEAD_SHA']
  resolved=run('resolve-published-evmd',['go','list','-m','-json','github.com/xitcoin-org/pos-chain/evmd@'+head])
  if resolved:
   resolved_text=(out/'resolve-published-evmd.log').read_text()
   version=json.JSONDecoder().raw_decode(resolved_text[resolved_text.index('{'):])[0]['Version']
   harness=root/'binary-qualification';harness.mkdir()
   shutil.copyfile(cwd/'go.mod',harness/'go.mod');shutil.copyfile(cwd/'go.sum',harness/'go.sum');cwd=harness
   run('prepare-versioned-binary',['go','mod','edit','-module=xitcoin.invalid/qualification','-require=github.com/xitcoin-org/pos-chain/evmd@'+version])
  else:version=None
 if module=='evmd' and version and run('build-symbols',['go','build','-mod=mod','-o',str(out/'xitcoind-symbols'),'github.com/xitcoin-org/pos-chain/evmd/cmd/evmd']):
  run('buildinfo',['go','version','-m',str(out/'xitcoind-symbols')])
  if '(devel)' in (out/'buildinfo.log').read_text():results['buildinfo']['exit_code']=1
  run('binary-symbols-scan',[scanner,'-mode=binary','-json',str(out/'xitcoind-symbols')])
  import hashlib
  (out/'xitcoind-symbols.sha256').write_text(hashlib.sha256((out/'xitcoind-symbols').read_bytes()).hexdigest()+'  xitcoind-symbols\n')
  # Keep the binary on the runner only; evidence contains build info, hash, and full scan.
  (out/'xitcoind-symbols').unlink()
# JSON scanner exit status alone is not a security disposition. Preserve every
# finding and its original trace, including modules hidden by fork coordinates.
scan_findings={}
for label in ['source-scan','binary-symbols-scan']:
 path=out/(label+'.log')
 if not path.exists(): continue
 raw=path.read_text(); findings=[]; decoder=json.JSONDecoder()
 try:
  while raw.strip():
   obj,end=decoder.raw_decode(raw.lstrip());raw=raw.lstrip()[end:]
   if 'finding' in obj: findings.append(obj['finding'])
 except json.JSONDecodeError as err:
  results[label]['parse_error']=str(err);results[label]['exit_code']=1
 scan_findings[label]={'ids':sorted({f['osv'] for f in findings}), 'findings':findings}
(out/'scanner-findings.json').write_text(json.dumps(scan_findings,indent=2))
failed=forbidden or raw_key_files or any(v['exit_code']!=0 or v['blocked'] for v in results.values())
(out/'results.json').write_text(json.dumps(results,indent=2))
(out/'qualification.json').write_text(json.dumps({'status':'QUALIFICATION_FAILED' if failed else 'TECHNICAL_CHECKS_COMPLETED_REVIEW_REQUIRED','security_acceptance':False,'review_expired':'2026-09-05','deployment':False},indent=2))
sys.exit(bool(failed))
