import { zClassSchema, zMakeErrorDataSchema } from '$lib/models';
import { processError } from '$lib/pocketbase';
import { redirect, type Actions, type ServerLoad } from '@sveltejs/kit';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

const zFormSchema = zClassSchema;
const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
const zSaveSchema = zClassSchema.omit({ attachments: true });

export const load: ServerLoad = async ({ locals, params }) => {
  const { id } = params;
  const { user, pb } = locals;

  // don't allow not logged in here
  if (!id || !user) throw redirect(303, '/');

  const klass = await pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zFormSchema.parse(r));

  if (klass.owner !== user.id && user.role !== 'editor') throw redirect(303, '/');

  const form = await superValidate(klass, zod(zFormSchema));

  return { klass, form };
};

export const actions: Actions = {
  save: async ({ locals, request, params }) => {
    const { pb } = locals;
    const { id } = params;

    const form = await superValidate(request, zod(zSaveSchema));
    if (!form.valid || !id) {
      return fail(400, { form });
    }

    try {
      await pb.collection('classes').update(id, form.data);
    } catch (e) {
      return processError(form, e, zErrorSchema);
    }

    return message(form, 'Salvo');
  },
  remove: async ({ locals, params }) => {
    const { pb } = locals;
    const { id } = params;
    if (!id) throw new Error('missing id');

    await pb.collection('classes').delete(id);

    throw redirect(303, '/');
  },
  removefile: async ({ locals, params, request }) => {
    const { pb } = locals;
    const { id } = params;
    if (!id) throw new Error('missing id');

    const data = await request.formData();
    const filename = data.get('attachment') as string;

    await pb.collection('classes').update(id, { 'attachments-': filename });
  },
};
