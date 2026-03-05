import api from '../api/client';

/**
 * Convert a base64 VAPID key to a Uint8Array for the Web Push API
 */
function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}

/**
 * Register for push notifications.
 * Call this after the user logs in.
 */
export async function registerPushNotifications() {
    // Check browser support
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
        console.log('Push notifications not supported in this browser');
        return false;
    }

    try {
        // Request notification permission
        const permission = await Notification.requestPermission();
        if (permission !== 'granted') {
            console.log('Notification permission denied');
            return false;
        }

        // Register service worker
        const registration = await navigator.serviceWorker.register('/sw.js');
        console.log('Service Worker registered');

        // Get VAPID public key from backend
        const { data } = await api.get('/push/vapid-key');
        const vapidKey = urlBase64ToUint8Array(data.public_key);

        // Get existing subscription
        let subscription = await registration.pushManager.getSubscription();

        // If subscription exists, check if VAPID key matches
        if (subscription) {
            const currentKey = subscription.options.applicationServerKey;
            const currentKeyUint8 = new Uint8Array(currentKey);

            // Compare keys
            let mismatch = false;
            if (currentKeyUint8.length !== vapidKey.length) {
                mismatch = true;
            } else {
                for (let i = 0; i < vapidKey.length; i++) {
                    if (currentKeyUint8[i] !== vapidKey[i]) {
                        mismatch = true;
                        break;
                    }
                }
            }

            if (mismatch) {
                console.log('VAPID key mismatch detected, re-subscribing...');
                await subscription.unsubscribe();
                subscription = null;
            }
        }

        if (!subscription) {
            subscription = await registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: vapidKey,
            });
            console.log('Push subscription created');
        }

        // Send subscription to backend
        const sub = subscription.toJSON();
        await api.post('/push/subscribe', {
            endpoint: sub.endpoint,
            keys: {
                p256dh: sub.keys.p256dh,
                auth: sub.keys.auth,
            },
        });

        console.log('Push notifications registered successfully');
        return true;
    } catch (error) {
        console.error('Failed to register push notifications:', error);
        return false;
    }
}

/**
 * Unregister from push notifications.
 */
export async function unregisterPushNotifications() {
    try {
        const registration = await navigator.serviceWorker.ready;
        const subscription = await registration.pushManager.getSubscription();

        if (subscription) {
            await api.delete('/push/unsubscribe', {
                data: { endpoint: subscription.endpoint },
            });
            await subscription.unsubscribe();
            console.log('Push notifications unregistered');
        }
    } catch (error) {
        console.error('Failed to unregister push notifications:', error);
    }
}
