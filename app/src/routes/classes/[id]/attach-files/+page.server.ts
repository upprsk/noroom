import { zClassSchema, zFileUploadSchema } from '$lib/models';
import { redirect, type Actions, type ServerLoad } from '@sveltejs/kit';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

const zFormSchema = zFileUploadSchema;

export const load: ServerLoad = async ({ locals, params }) => {
  const { id } = params;
  const { user, pb } = locals;

  // don't allow not logged in here
  if (!id || !user || user.role !== 'editor') throw redirect(303, '/');

  const klass = await pb
    .collection('classes')
    .getOne(id, { fetch })
    .then((r) => zClassSchema.parse(r));

  const form = await superValidate(zod(zFormSchema));

  return { klass, form };
};

export const actions: Actions = {
  default: async ({ params, locals, request }) => {
    const { pb } = locals;
    const { id } = params;

    const form = await superValidate(request, zod(zFormSchema));

    if (!form.valid || !id) {
      return fail(400, { form });
    }

    console.log('form:', form.data.attachments);

    const r = await pb.collection('classes').update(id, form.data);
    console.log('response:', r);

    return message(form, `${form.data.attachments.length} arquivos anexados`);
  },
};
