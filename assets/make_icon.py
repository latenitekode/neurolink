"""Generate Neurolink's portable Windows icon (no binary files committed)."""
from pathlib import Path
from PIL import Image, ImageDraw

SIZE = 256
im = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
d = ImageDraw.Draw(im)
d.rounded_rectangle((8, 8, 248, 248), radius=55, fill="#101829", outline="#36E0CC", width=7)
# Radar rings and crosshair, evoking Crypto's surveillance drone.
for margin in (44, 70):
    d.ellipse((margin, margin, SIZE-margin, SIZE-margin), outline="#267F91", width=4)
d.line((128, 34, 128, 222), fill="#247889", width=3)
d.line((34, 128, 222, 128), fill="#247889", width=3)
# Stylized N lightning glyph, bold enough to remain legible at 16px.
d.line([(78, 180), (78, 76), (174, 180), (174, 76)], fill="#3CF4DA", width=19, joint="curve")
d.ellipse((166, 68, 182, 84), fill="#E8FFFB")
out = Path("assets/neurolink.ico")
out.parent.mkdir(parents=True, exist_ok=True)
im.save(out, format="ICO", sizes=[(16,16),(24,24),(32,32),(48,48),(64,64),(128,128),(256,256)])
print(f"Created {out}")
