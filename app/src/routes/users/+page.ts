import { zUserSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

const zUsersArraySchema = zUserSchema.array();

export const load: Load = async ({ fetch }) => {
  const user = get(currentUser);

  if (!user) throw redirect(303, '/');

  const usersP = pb
    .collection('users')
    .getFullList({ fetch })
    .then((r) => zUsersArraySchema.parse(r));

  const [users] = await Promise.all([usersP]);

  return { users };
};
