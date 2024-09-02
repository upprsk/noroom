import { zClassWithPresenceSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load: Load = async ({ params, fetch }) => {
  const { id } = params;

  if (!id || !get(currentUser)) throw redirect(303, '/');

  const classP = pb
    .collection('classes')
    .getOne(id, {
      fetch,
      expand: 'classPresenceEntries_via_class,classPresenceEntries_via_class.user',
    })
    .then((r) => zClassWithPresenceSchema.parse(r));

  const [klass] = await Promise.all([classP]);

  return { klass };
};
