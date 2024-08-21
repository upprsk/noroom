import { pb } from './pocketbase';

export const sendTracking = async (
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
