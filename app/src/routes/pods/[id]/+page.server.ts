import { zPodWithServerSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';

export const load: ServerLoad = async ({ params, locals, fetch }) => {
  const { id } = params;
  const { pb, user } = locals;

  if (!user) throw redirect(303, '/');
  if (!id) throw redirect(303, '/pods');

  const podP = pb
    .collection('pods')
    .getOne(id, { fetch, expand: 'server' })
    .then((r) => zPodWithServerSchema.parse(r));

  const [pod] = await Promise.all([podP]);

  return { pod };
};
