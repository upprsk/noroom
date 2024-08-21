<script lang="ts">
  export let name: string;
  export let errors: string[] | undefined;
  export let value: string;
  export let constraints: { [key: string]: unknown } | undefined;

  export let disabled = false;
  export let placeholder = name;
</script>

<label class="form-control w-full max-w-xs">
  <div class="label">
    <span class="label-text"><slot /></span>
  </div>
  <input
    type="email"
    class="input input-bordered w-full max-w-xs"
    aria-invalid={errors ? 'true' : undefined}
    bind:value
    {name}
    {disabled}
    {placeholder}
    {...constraints}
  />
  <div class="label">
    {#if errors}<span class="label-text-alt text-error">{errors}</span>{/if}
  </div>
</label>
