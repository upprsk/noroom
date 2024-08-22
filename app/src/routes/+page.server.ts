import { zClassSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import type { ServerLoad } from '@sveltejs/kit';

const zClassesArraySchema = zClassSchema.array();

export const load: ServerLoad = async ({ locals, fetch }) => {
  const { pb } = locals;

  const classesP = pb
    .collection('classes')
    .getFullList({ fetch })
    .then((r) => zClassesArraySchema.parse(r));

  const [classes] = await Promise.all([classesP]);

  return { classes };
};
