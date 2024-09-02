import { zPodSchema, zPodServerSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

const zPodServerArraySchema = zPodServerSchema.array();
const zPodArraySchema = zPodSchema.array();

export const load: Load = async ({ fetch, depends }) => {
  const user = get(currentUser);

  if (!user) throw redirect(303, '/');

  const podServersP = pb
    .collection('podServers')
    .getFullList({ fetch })
    .then((r) => zPodServerArraySchema.parse(r));

  const podsP = pb
    .collection('pods')
    .getFullList({ fetch })
    .then((r) => zPodArraySchema.parse(r));

  const [podServers, pods] = await Promise.all([podServersP, podsP]);

  depends('app:podServers');

  return { user, podServers, pods };
};
