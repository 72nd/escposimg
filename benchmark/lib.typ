#let tpl(
  height: auto,
  font_size: 8pt,
  right_margin: 7.5mm,
  body,
) = {
  set page(width: 79.5mm, height: height, margin: (
    top: 0mm,
    left: 0mm,
    right: right_margin,
    bottom: 0mm,
  ))

  set par(
    spacing: 0.9em,
    leading: 0.4em,
  )

  show heading.where(level: 2): set text(size: 0.9em)
  show heading.where(level: 2): set block(above: .9em)

  set text(
    size: font_size,
    font: "Arial",
    hyphenate: true,
  )

  show raw: set text(font: "Courier New", size: 10pt)

  show heading.where(level: 1): set text(size: 0.6em)

  body
}
