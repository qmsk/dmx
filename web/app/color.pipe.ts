import { Pipe, PipeTransform } from '@angular/core';

import { Color, Value } from './types';

import * as _ from 'lodash';

@Pipe({
    name: 'color',
})
export class ColorPipe implements PipeTransform {
  transform(color: Color): string {
    return this.hexColor(color);
  }

  private hexValue(value: Value): string {
    return _.padStart(Math.trunc(value * 255).toString(16), 2, '0');
  }

  private hexColor(color: Color): string {
    return "#" + this.hexValue(color.Red) + this.hexValue(color.Green) + this.hexValue(color.Blue);
  }
}
