<script lang="ts">
  import { applyAction, enhance } from '$app/forms';
  import { pb } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
  import { clrLastTrackedTime } from '$lib/tracking';
  import BasicAvatar from './BasicAvatar.svelte';
</script>

<div class="navbar-end gap-2">
  {#if $currentUser}
    <form
      method="POST"
      action="/logout"
      use:enhance={() =>
        async ({ result }) => {
          pb.authStore.clear();
          clrLastTrackedTime();
          await applyAction(result);
        }}
    >
      <button class="btn btn-ghost btn-sm">Logout</button>
    </form>

    <BasicAvatar href="/account" user={$currentUser} />
  {:else}
    <a href="/login" class="btn btn-ghost btn-sm">Sign-In</a>
    <a href="/register" class="btn btn-accent btn-sm">Sign-Up</a>
  {/if}
</div>
