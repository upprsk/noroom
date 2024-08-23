<script lang="ts">
  import { placeholderLetters } from '$lib';
  import type { UserModel } from '$lib/models';
  import { getFileUrl } from '$lib/pocketbase';

  export let user: UserModel;
  export let href: string;
</script>

{#if user.avatar !== ''}
  <a {href} class="avatar">
    <div class="w-10 rounded-full">
      <img alt="avatar" src={getFileUrl(user, user.avatar, { thumb: '64x64' })} />
    </div>
  </a>
{:else}
  <a {href} class="avatar placeholder">
    <div class="bg-neutral text-neutral-content w-10 rounded-full">
      <span>{placeholderLetters(user.name)}</span>
    </div>
  </a>
{/if}
