import type { UserModel } from '$lib/models';
import { writable } from 'svelte/store';

export const currentUser = writable<UserModel | null>();
