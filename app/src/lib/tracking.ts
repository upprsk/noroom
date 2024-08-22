import type Client from 'pocketbase';

export const sendTracking = async (
  pb: Client,
  userid: string | undefined,
  fingerprint: string,
  deviceData: unknown,
) => {
  try {
    await pb.send('/api/noroom/tracking', {
      method: 'POST',
      body: { userid, fingerprint, deviceData },
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
export const clrLastTrackedTime = () => localStorage.removeItem(TRACKING_KEY);
