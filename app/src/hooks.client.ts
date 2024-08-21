import { zUserSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';

pb.authStore.loadFromCookie(document.cookie);
pb.authStore.onChange(() => {
  const model = zUserSchema.nullable().parse(pb.authStore.model);

  currentUser.set(model);
  document.cookie = pb.authStore.exportToCookie({ httpOnly: false });
}, true);
