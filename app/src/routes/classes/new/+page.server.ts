import { zClassSchema, zMakeErrorDataSchema } from '$lib/models';
import { processError } from '$lib/pocketbase';
import { redirect, type Actions, type ServerLoad } from '@sveltejs/kit';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

const zFormSchema = zClassSchema;
const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());

export const load: ServerLoad = async ({ locals }) => {
  // don't allow not logged in here
  if (!locals.user || locals.user.role !== 'editor') throw redirect(303, '/');

  const form = await superValidate(zod(zFormSchema));

  return { form };
};

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const form = await superValidate(request, zod(zFormSchema));

    if (!form.valid) {
      // Again, return { form } and things will just work.
      return fail(400, { form });
    }

    try {
      await locals.pb.collection('classes').create(form.data);
    } catch (e) {
      return processError(form, e, zErrorSchema);
    }

    redirect(303, '/');
  },
};
