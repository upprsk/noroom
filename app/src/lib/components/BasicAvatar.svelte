<script lang="ts">
  import type { UserModel } from '$lib/models';
  import { getImageUrl } from '$lib/pocketbase';

  export let user: UserModel;

  const placeholderLetters = (name: string) =>
    name
      .split(' ')
      .filter((w) => w !== '')
      .map((w) => w[0])
      .join('')
      .toUpperCase();
</script>

{#if user.avatar !== ''}
  <div class="avatar">
    <div class="w-12 rounded-full">
      <img alt="avatar" src={getImageUrl(user, user.avatar, { thumb: '64x64' })} />
    </div>
  </div>
{:else}
  <div class="avatar placeholder">
    <div class="bg-neutral text-neutral-content w-12 rounded-full">
      <span>{placeholderLetters(user.name)}</span>
    </div>
  </div>
{/if}
