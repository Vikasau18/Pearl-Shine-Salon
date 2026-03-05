// Service Worker for Push Notifications
// This file must be served from the root of the domain

self.addEventListener('push', function (event) {
    if (!event.data) return;

    let data;
    try {
        data = event.data.json();
    } catch (e) {
        data = {
            title: 'Saloon',
            body: event.data.text(),
        };
    }

    const options = {
        body: data.body || 'You have a new notification',
        icon: data.icon || '/vite.svg',
        badge: '/vite.svg',
        vibrate: [200, 100, 200],
        data: {
            url: data.url || '/',
        },
        actions: [
            { action: 'open', title: 'Open' },
            { action: 'dismiss', title: 'Dismiss' },
        ],
    };

    event.waitUntil(
        self.registration.showNotification(data.title || 'Saloon', options)
    );
});

// Handle notification click
self.addEventListener('notificationclick', function (event) {
    event.notification.close();

    if (event.action === 'dismiss') return;

    const urlToOpen = event.notification.data?.url || '/';

    event.waitUntil(
        clients.matchAll({ type: 'window', includeUncontrolled: true }).then(function (clientList) {
            // If there's already an open window, focus it
            for (const client of clientList) {
                if (client.url.includes(self.location.origin) && 'focus' in client) {
                    client.navigate(urlToOpen);
                    return client.focus();
                }
            }
            // Otherwise open a new window
            return clients.openWindow(urlToOpen);
        })
    );
});

// Activate immediately
self.addEventListener('activate', function (event) {
    event.waitUntil(self.clients.claim());
});
