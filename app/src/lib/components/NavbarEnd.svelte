<script lang="ts">
  import { applyAction, enhance } from '$app/forms';
  import { pb } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
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
          await applyAction(result);
        }}
    >
      <button class="btn btn-ghost btn-sm">Logout</button>
    </form>

    <BasicAvatar href="/account" user={$currentUser} />
  {:else}
    <a href="/login" class="btn-sm btn btn-ghost">Sign-In</a>
    <a href="/register" class="btn-sm btn btn-accent">Sign-Up</a>
  {/if}
</div>
