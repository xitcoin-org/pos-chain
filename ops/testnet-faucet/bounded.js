'use strict';
const { spawn } = require('node:child_process');

function readBody(stream, limit = 2048, timeout = 5000) {
  return new Promise((resolve, reject) => {
    let bytes = 0;
    const chunks = [];
    const finish = (err) => {
      clearTimeout(timer);
      stream.removeListener('data', data);
      stream.removeListener('end', end);
      stream.removeListener('error', error);
      stream.removeListener('aborted', aborted);
      if (err) { stream.pause(); reject(err); } else resolve(Buffer.concat(chunks).toString('utf8'));
    };
    const data = (chunk) => { bytes += chunk.length; if (bytes > limit) finish(new Error('body_limit')); else chunks.push(Buffer.from(chunk)); };
    const end = () => finish();
    const error = (err) => finish(err);
    const aborted = () => finish(new Error('body_aborted'));
    const timer = setTimeout(() => finish(new Error('body_timeout')), timeout);
    stream.on('data', data).once('end', end).once('error', error).once('aborted', aborted);
  });
}

function runBounded(command, args, options, { spawnChild = spawn, timeout = 30000, limit = 65536 } = {}) {
  return new Promise((resolve, reject) => {
    const child = spawnChild(command, args, options);
    let output = '', bytes = 0, done = false;
    const finish = (err) => {
      if (done) return;
      done = true;
      clearTimeout(timer);
      if (err) { child.kill('SIGKILL'); reject(err); } else resolve(output);
    };
    const timer = setTimeout(() => finish(new Error('child_timeout_unknown')), timeout);
    child.stdout.on('data', (chunk) => { if (done) return; bytes += chunk.length; if (bytes > limit) finish(new Error('child_output_unknown')); else output += chunk.toString(); });
    child.stderr.on('data', (chunk) => { if (done) return; bytes += chunk.length; if (bytes > limit) finish(new Error('child_output_unknown')); });
    child.once('error', () => finish(new Error('child_error_unknown')));
    child.once('close', (code) => finish(code === 0 ? null : new Error('child_exit_unknown')));
  });
}
module.exports = { readBody, runBounded };
