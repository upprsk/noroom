import { zMakeErrorDataSchema, zRegisterSchema } from '$lib/models';
import { processError } from '$lib/pocketbase';
import { fail, redirect, type Actions, type ServerLoad } from '@sveltejs/kit';
import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

const zFormSchema = zRegisterSchema;
const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());

export const load: ServerLoad = async ({ locals }) => {
  // don't allow logged in users here
  if (locals.user) throw redirect(303, '/account');

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
      await locals.pb.collection('users').create({ ...form.data, role: 'student' });
      await locals.pb.collection('users').authWithPassword(form.data.email, form.data.password);
    } catch (e) {
      return processError(form, e, zErrorSchema);
    }

    redirect(303, '/');
  },
};
