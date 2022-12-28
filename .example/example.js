const spawn = require('child_process').spawn;

let time;
let processData = (s) => {
  let data = s.replace(/}\\n/g, '},');
  if (data.endsWith(',')) data = data.slice(0, -1).trim();
  let messages = JSON.parse(`[${data}]`);
  for (let message of messages) {
    if (message.e === 'schedule:done') {
      console.log(Date.now() - time);
    }
    console.log('[%s]:', message.e, message.e.includes('ready') ? new Date() : message.d);
  }
};

const go = spawn('./scheduler.exe');

let sendData = (e, d) => go.stdin.write(JSON.stringify({ e, d }) + '\n');

go.stdout.setEncoding('utf8');
go.stderr.setEncoding('utf8');

go.stderr.on('error', console.error);
go.stderr.on('data', console.error);
go.stdout.on('data', processData);

sendData('schedule:add', {
  id: '1234', // id unique (this id is use to stored in db)
  expires_at: 10, // seconds (Date time handle in scheduler.exe)
  content: 'Expire 1234 in 10s' // maybe use JSON.stringify to stored data
});

sendData('schedule:exists', '1234'); // return event: schedule:exists, d: true

// sendData('schedule:delete', '1234') // remove

sendData('schedule:add', {
  id: '12345',
  expires_at: 20,
  content: 'Expire 1235 in 20s'
});
time = Date.now();

process.on('beforeExit', () => {
  sendData('exit', {});
  go.kill();
});