import { zEndDeviceSchema } from '$lib/models';
import type { ServerLoad } from '@sveltejs/kit';

const zDevices = zEndDeviceSchema.array();

export const load: ServerLoad = async ({ locals, fetch }) => {
  const { pb } = locals;

  const devicesP = pb
    .collection('endDevices')
    .getFullList({ fetch })
    .then((r) => zDevices.parseAsync(r));

  const [devices] = await Promise.all([devicesP]);

  return { devices };
};
