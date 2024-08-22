import { zClassSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

export const ssr = false;

export const load: ServerLoad = async ({ locals, params, fetch }) => {
  const { id } = params;
  const { pb, user } = locals;

  if (!id || !user) throw redirect(303, '/');

  const classP = pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zClassSchema.parse(r));

  const [klass] = await Promise.all([classP]);

  return { klass };
};
