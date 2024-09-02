import { zEndDeviceSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

const zDevices = zEndDeviceSchema.array();

export const load: Load = async ({ fetch }) => {
  const user = get(currentUser);

  if (!user) redirect(303, '/');

  const devicesP = pb
    .collection('endDevices')
    .getFullList({ fetch })
    .then((r) => zDevices.parseAsync(r));

  const [devices] = await Promise.all([devicesP]);

  return { user, devices };
};
