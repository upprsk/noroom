import { zEndDeviceSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

const zDevices = zEndDeviceSchema.array();

export const load: ServerLoad = async ({ locals, fetch }) => {
  const { pb, user } = locals;

  if (!user) redirect(303, '/');

  const devicesP = pb
    .collection('endDevices')
    .getFullList({ fetch })
    .then((r) => zDevices.parseAsync(r));

  const [devices] = await Promise.all([devicesP]);

  return { devices };
};
