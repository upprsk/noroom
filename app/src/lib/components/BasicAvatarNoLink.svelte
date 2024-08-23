<script lang="ts">
  import { placeholderLetters } from '$lib';
  import type { UserModel } from '$lib/models';
  import { getFileUrl } from '$lib/pocketbase';

  export let user: UserModel;
</script>

{#if user.avatar !== ''}
  <div class="avatar">
    <div class="w-10 rounded-full bg-base-100">
      <img alt="avatar" src={getFileUrl(user, user.avatar, { thumb: '64x64' })} />
    </div>
  </div>
{:else}
  <div class="avatar placeholder">
    <div class="w-10 rounded-full bg-neutral text-neutral-content">
      <span>{placeholderLetters(user.name)}</span>
    </div>
  </div>
{/if}
