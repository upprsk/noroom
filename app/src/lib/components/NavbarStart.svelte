<script>
  import NavbarMenu from './NavbarMenu.svelte';
</script>

<div class="navbar-start">
  <div class="dropdown">
    <div tabindex="0" role="button" class="btn btn-ghost lg:hidden">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 12h8m-8 6h16"
        />
      </svg>
    </div>
    <NavbarMenu
      tabindex="0"
      class="menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
    />
  </div>
  <a href="/" class="btn btn-ghost text-xl">No Room</a>
</div>
