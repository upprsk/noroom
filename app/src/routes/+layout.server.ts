import type { UserModel } from '$lib/models';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
  return {
    user: locals.user as UserModel,
  };
};
