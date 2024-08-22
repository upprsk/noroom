import type Client from 'pocketbase';

export const sendTracking = async (
  pb: Client,
  userid: string | undefined,
  fingerprint: string,
  deviceData: unknown,
) => {
  // get current ip address information
  let locationData;
  try {
    locationData = await fetch(`http://ip-api.com/json`).then((r) => r.json());
  } catch (e) {
    console.error('Failed to get location:', e);
  }

  try {
    await pb.send('/api/noroom/tracking', {
      method: 'POST',
      body: { userid, fingerprint, deviceData, locationData },
    });
  } catch (e) {
    console.error(e);
  }
};

const TRACKING_KEY = 'noroom:tracking_time';
export const getLastTrackedTime = () =>
  new Date(parseInt(localStorage.getItem(TRACKING_KEY) ?? ''));
export const setLastTrackedTime = () =>
  localStorage.setItem(TRACKING_KEY, new Date().getTime().toString());
