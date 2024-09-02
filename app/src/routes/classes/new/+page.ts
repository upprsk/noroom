import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load: Load = async () => {
  const user = get(currentUser);

  // don't allow not logged in here
  if (!user || user.role !== 'editor') throw redirect(303, '/');
};
