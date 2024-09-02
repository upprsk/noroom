<script lang="ts">
  import { goto } from '$app/navigation';
  import { pb } from '$lib/pocketbase';
  import { currentUser, shouldTrack } from '$lib/stores/user';
  import { clrLastTrackedTime } from '$lib/tracking';
  import BasicAvatar from './BasicAvatar.svelte';
</script>

<div class="navbar-end gap-2">
  {#if $currentUser}
    <button
      type="button"
      class="btn btn-ghost btn-sm"
      on:click={() => {
        pb.authStore.clear();
        clrLastTrackedTime();

        goto('/');
        shouldTrack.set(true);
      }}
    >
      Logout
    </button>

    <BasicAvatar href="/account" user={$currentUser} />
  {:else}
    <a href="/login" class="btn btn-ghost btn-sm">Sign-In</a>
    <a href="/register" class="btn btn-accent btn-sm">Sign-Up</a>
  {/if}
</div>
