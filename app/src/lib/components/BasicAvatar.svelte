<script lang="ts">
  import type { UserModel } from '$lib/models';
  import { getImageUrl } from '$lib/pocketbase';

  export let user: UserModel;
  export let href: string;

  const placeholderLetters = (name: string) =>
    name
      .split(' ')
      .filter((w) => w !== '')
      .map((w) => w[0])
      .join('')
      .toUpperCase();
</script>

{#if user.avatar !== ''}
  <a {href} class="avatar">
    <div class="w-10 rounded-full">
      <img alt="avatar" src={getImageUrl(user, user.avatar, { thumb: '64x64' })} />
    </div>
  </a>
{:else}
  <a {href} class="avatar placeholder">
    <div class="bg-neutral text-neutral-content w-10 rounded-full">
      <span>{placeholderLetters(user.name)}</span>
    </div>
  </a>
{/if}
