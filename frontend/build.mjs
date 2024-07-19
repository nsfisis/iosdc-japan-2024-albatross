import esbuild from 'esbuild';

await esbuild.build({
  entryPoints: ['src/game.jsx', 'src/watch.jsx'],
  outdir: 'dist/js',
  bundle: true,
  // minify: true,
  minify: false,
  sourcemap: true,
  platform: 'browser',
  format: 'esm',
  jsx: 'automatic',
});
