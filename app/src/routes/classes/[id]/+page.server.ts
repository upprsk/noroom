import { zClassSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

export const load: ServerLoad = async ({ locals, params, fetch }) => {
  const { id } = params;
  const { pb } = locals;

  if (!id) throw redirect(303, '/');

  const classP = pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zClassSchema.parse(r));

  const [klass] = await Promise.all([classP]);

  return { klass };
};
