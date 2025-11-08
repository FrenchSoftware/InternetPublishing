// Hot reload script - only active in development
(function() {
  if (location.hostname !== 'localhost' && location.hostname !== '127.0.0.1') {
    return; // Only run on localhost
  }

  let ws;
  
  function connect() {
    ws = new WebSocket(`ws://${location.host}/__hotreload`);
    
    ws.onopen = () => {
      console.log('[Hot Reload] Connected');
    };
    
    ws.onmessage = (event) => {
      if (event.data === 'reload') {
        console.log('[Hot Reload] Reloading page...');
        location.reload();
      }
    };
    
    ws.onclose = () => {
      console.log('[Hot Reload] Disconnected. Retrying in 1s...');
      setTimeout(connect, 1000);
    };
    
    ws.onerror = (error) => {
      console.log('[Hot Reload] Error:', error);
      ws.close();
    };
  }
  
  connect();
})();
