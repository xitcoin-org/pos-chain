import hashlib,json,os,pathlib,shutil,signal,subprocess,sys,time
name=sys.argv[1];root=pathlib.Path.cwd();out=root/'qualification-evidence';out.mkdir(exist_ok=True)
base={'cosmos-sdk':'dedeb7c80a91c47ae83f5352e29c3dd34e4a3fc6','go-ethereum':'d99d6fa2c8d98b7cd653de4a9386d2da3db8f25c'}[name]
src=root/'qualification-source';start=time.monotonic();minimum=shutil.disk_usage(root).free

def run(label,args,cwd=src):
 global minimum
 with (out/(label+'.log')).open('w') as log:
  child_env=os.environ.copy()
  if cwd==src:child_env.update(GITHUB_SHA=base,GITHUB_REPOSITORY='cosmos/'+name,GITHUB_REF_TYPE='',GITHUB_HEAD_REF='')
  p=subprocess.Popen(args,cwd=cwd,env=child_env,stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
  reason=None;t=time.monotonic()
  while p.poll() is None:
   free=shutil.disk_usage(root).free;minimum=min(minimum,free)
   mem=int(next(l.split()[1] for l in pathlib.Path('/proc/meminfo').read_text().splitlines() if l.startswith('MemAvailable:')))*1024
   if free<5.25*1024**3 or mem<768*1024**2 or time.monotonic()-t>7200:
    reason='resource guard or 7200s stage limit';os.killpg(p.pid,signal.SIGTERM)
    try:p.wait(timeout=10)
    except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
    break
   time.sleep(.5)
 result={'command':args,'exit_code':p.returncode,'blocked':reason,'minimum_free_bytes':minimum,'seconds':round(time.monotonic()-t,2)}
 (out/(label+'.json')).write_text(json.dumps(result,indent=2));print(label,json.dumps(result),flush=True)
 if p.returncode or reason:
  print((out/(label+'.log')).read_text()[-12000:]);sys.exit(1)

run('clone',['git','clone','--filter=blob:none','--no-checkout','--depth=1','https://github.com/cosmos/'+name+'.git',str(src)],root)
run('fetch',['git','fetch','--depth=1','origin',base]);run('checkout',['git','checkout','--detach',base])
patch=root/'docs/fork-qualification-20260913'/(name+'.patch')
(out/'input.json').write_text(json.dumps({'name':name,'base':base,'patch_sha256':hashlib.sha256(patch.read_bytes()).hexdigest(),'qualification_commit':os.getenv('GITHUB_SHA')},indent=2))
run('patch',['git','apply','--unidiff-zero',str(patch)])
if name=='cosmos-sdk':
 for label,args in [('tidy-all',['make','tidy-all','VERSION_RAW=v0.54.4']),('build',['make','build','VERSION_RAW=v0.54.4']),('lint',['make','lint','VERSION_RAW=v0.54.4']),('test-unit',['make','test-unit','VERSION_RAW=v0.54.4'])]:run(label,args)
else:
 run('gofmt',['gofmt','-w','p2p/nat/stun.go','p2p/nat/stun_local_test.go'])
 run('goimports',['go','run','golang.org/x/tools/cmd/goimports@v0.41.0','-w','p2p/nat/stun.go','p2p/nat/stun_local_test.go'])
 for label,args in [('build',['make','all']),('test',['go','run','./build/ci.go','test']),('lint',['go','run','./build/ci.go','lint']),('check-generate',['go','run','./build/ci.go','check_generate']),('check-baddeps',['go','run','./build/ci.go','check_baddeps'])]:run(label,args)
(out/'complete.json').write_text(json.dumps({'status':'REQUIRED_UPSTREAM_CHECKS_PASSED','base':base,'patch_sha256':hashlib.sha256(patch.read_bytes()).hexdigest(),'minimum_free_bytes':minimum},indent=2))
