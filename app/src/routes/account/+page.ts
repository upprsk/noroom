import { zEndDeviceSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import type { PageLoad } from './$types';

const zDevices = zEndDeviceSchema.array();

export const load: PageLoad = async ({ fetch }) => {
  const devicesP = pb
    .collection('endDevices')
    .getFullList({ fetch })
    .then((r) => zDevices.parseAsync(r));

  const [devices] = await Promise.all([devicesP]);

  return { devices };
};
