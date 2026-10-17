// Workspace authentication needs a live navigation, including when an older
// Music PWA has cached its SPA shell. Leave its asset/audio caches in place.
self.addEventListener('fetch', event => {
  const url = new URL(event.request.url);
  if (url.origin !== self.location.origin || !url.pathname.startsWith('/music/') || event.request.mode !== 'navigate') return;
  event.stopImmediatePropagation();
  event.respondWith(fetch(event.request).catch(async () => {
    if (!url.pathname.startsWith('/music/auth/')) {
      const shell = await caches.match(new URL('index.html', self.registration.scope).href, {ignoreSearch: true});
      if (shell) return shell;
    }
    return new Response('Music needs a connection to Workspace. Reconnect and open Music again.', {status: 503, headers: {'Content-Type': 'text/plain; charset=utf-8'}});
  }));
});
self.addEventListener('install', event => event.waitUntil(self.skipWaiting()));
self.addEventListener('activate', event => event.waitUntil(self.clients.claim()));
