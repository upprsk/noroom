import { zPodWithServerSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load: Load = async ({ params, fetch }) => {
  const { id } = params;

  const user = get(currentUser);

  if (!user) throw redirect(303, '/');
  if (!id) throw redirect(303, '/pods');

  const podP = pb
    .collection('pods')
    .getOne(id, { fetch, expand: 'server' })
    .then((r) => zPodWithServerSchema.parse(r));

  const [pod] = await Promise.all([podP]);

  return { pod };
};
