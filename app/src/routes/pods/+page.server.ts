import { zPodSchema, zPodServerSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

const zPodServerArraySchema = zPodServerSchema.array();
const zPodArraySchema = zPodSchema.array();

export const load: ServerLoad = async ({ locals, fetch, depends }) => {
  const { pb, user } = locals;

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

  return { podServers, pods };
};
