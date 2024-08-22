import { zFileUploadSchema } from '$lib/models';
import { zClassSchema, zMakeErrorDataSchema } from '$lib/models';
import { processError } from '$lib/pocketbase';
import { redirect, type Actions, type ServerLoad } from '@sveltejs/kit';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

const zFormSchema = zClassSchema;
const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());

export const load: ServerLoad = async ({ locals, params }) => {
  const { id } = params;
  const { user, pb } = locals;

  // don't allow not logged in here
  if (!id || !user || user.role !== 'editor') throw redirect(303, '/');

  const klass = await pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zFormSchema.parse(r));

  const form = await superValidate(klass, zod(zFormSchema));

  return { klass, form };
};

export const actions: Actions = {
  save: async ({ locals, request, params }) => {
    const { id } = params;

    const form = await superValidate(request, zod(zFormSchema));
    if (!form.valid || !id) {
      return fail(400, { form });
    }

    try {
      await locals.pb.collection('classes').update(id, form.data);
    } catch (e) {
      return processError(form, e, zErrorSchema);
    }

    return message(form, 'Salvo');
  },
  remove: async ({ locals, params }) => {
    const { id } = params;
    if (!id) throw new Error('missing id');

    await locals.pb.collection('classes').delete(id);

    throw redirect(303, '/');
  },
};
