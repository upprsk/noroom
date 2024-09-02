import { redirect, type Load } from '@sveltejs/kit';
import { pb } from '$lib/pocketbase';
import { get } from 'svelte/store';
import { currentUser } from '$lib/stores/user';
import { zPodServerArraySchema } from './models';

export const load: Load = async ({ fetch }) => {
  const user = get(currentUser);

  // don't allow not logged in here
  if (!user) throw redirect(303, '/');

  const podServersP = pb
    .collection('podServers')
    .getFullList({ fetch })
    .then((r) => zPodServerArraySchema.parse(r));

  const [podServers] = await Promise.all([podServersP]);

  return { user, podServers };
};
