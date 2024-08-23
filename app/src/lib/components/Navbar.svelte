<script lang="ts">
  import NavbarEnd from './NavbarEnd.svelte';
  import NavbarMenu from './NavbarMenu.svelte';
  import NavbarStart from './NavbarStart.svelte';
</script>

<div class="navbar bg-base-100 shadow-md z-10">
  <NavbarStart />

  <div class="navbar-center hidden lg:flex">
    <NavbarMenu class="menu-horizontal px-1" />
  </div>

  <NavbarEnd />
</div>
