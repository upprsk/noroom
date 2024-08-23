// place files you want to import through the `$lib` alias in this folder.

export const placeholderLetters = (name: string) =>
  name
    .split(' ')
    .filter((w) => w !== '')
    .map((w) => w[0])
    .join('')
    .toUpperCase();
