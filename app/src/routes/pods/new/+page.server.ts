import { zMakeErrorDataSchema, zPodSchema, zPodServerSchema } from '$lib/models';
import { redirect, type ServerLoad } from '@sveltejs/kit';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import type { Actions } from './$types';
import { processError } from '$lib/pocketbase';

const zFormSchema = zPodSchema;
const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
const zPodServerArraySchema = zPodServerSchema.array();

export const load: ServerLoad = async ({ locals }) => {
  const { user, pb } = locals;

  // don't allow not logged in here
  if (!user) throw redirect(303, '/');

  const podServersP = pb
    .collection('podServers')
    .getFullList({ fetch })
    .then((r) => zPodServerArraySchema.parse(r));

  const [podServers] = await Promise.all([podServersP]);
  const form = await superValidate(zod(zFormSchema));

  return { podServers, form };
};

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const { pb } = locals;

    const form = await superValidate(request, zod(zFormSchema));

    if (!form.valid) {
      // Again, return { form } and things will just work.
      return fail(400, { form });
    }

    console.log(form.data);

    let res;
    try {
      res = await pb.collection('pods').create(form.data);
    } catch (e) {
      return processError(form, e, zErrorSchema);
    }

    throw redirect(303, `/pods/${res.id}`);
  },
};
