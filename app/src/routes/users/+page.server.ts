import { zUserSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

const zUsersArraySchema = zUserSchema.array();

export const load: ServerLoad = async ({ locals, fetch }) => {
  const { user, pb } = locals;

  if (!user) throw redirect(303, '/');

  const usersP = pb
    .collection('users')
    .getFullList({ fetch })
    .then((r) => zUsersArraySchema.parse(r));

  const [users] = await Promise.all([usersP]);

  return { users };
};
