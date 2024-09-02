import { zClassSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { type Load, redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load: Load = async ({ params, fetch }) => {
  const { id } = params;

  const user = get(currentUser);

  // don't allow not logged in here
  if (!id || !user) throw redirect(303, '/');

  const klass = await pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zClassSchema.parse(r));

  if (klass.owner !== user.id && user.role !== 'editor') throw redirect(303, '/');

  return { klass };
};
