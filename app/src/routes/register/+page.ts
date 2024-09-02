import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';

export const load: Load = async () => {
  const user = get(currentUser);

  // don't allow logged in users here
  if (user) throw redirect(303, '/account');
};
