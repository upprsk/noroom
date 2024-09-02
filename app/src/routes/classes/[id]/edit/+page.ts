import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { type Load, redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';
import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { zFormSchema } from './models';

export const load: Load = async ({ params, fetch, depends }) => {
  const { id } = params;
  const user = get(currentUser);

  // don't allow not logged in here
  if (!id || !user) throw redirect(303, '/');

  const klass = await pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zFormSchema.parse(r));

  if (klass.owner !== user.id && user.role !== 'editor') throw redirect(303, '/');

  const form = await superValidate(klass, zod(zFormSchema));

  depends('app:klass');

  return { klass, form };
};
