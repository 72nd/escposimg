#!/usr/bin/env python3
"""
ESC/POS Image Printer Script
Prints an image to a network printer using ESC/POS commands.
"""

import argparse
import sys
from pathlib import Path
from escpos.printer import Network
from PIL import Image


def print_image_to_network_printer(image_path: str, printer_ip: str, printer_port: int = 9100, profile: str = "TM-T20II", max_width: int = 576):
    """
    Print an image to a network printer using ESC/POS.
    
    Args:
        image_path: Path to the image file
        printer_ip: IP address of the printer
        printer_port: Port of the printer (default: 9100)
        profile: Printer profile (default: TM-T20II for TM-T20III compatibility)
        max_width: Maximum width in pixels for the image (default: 576 for 80mm paper)
    """
    try:
        # Check if image file exists
        if not Path(image_path).exists():
            print(f"Error: Image file '{image_path}' not found")
            return False
        
        # Open and prepare image
        print(f"Opening image: {image_path}")
        img = Image.open(image_path)
        
        # Convert to RGB if necessary
        if img.mode != 'RGB':
            img = img.convert('RGB')
        
        # Scale image to full printer width
        print(f"Original image size: {img.width}x{img.height}")
        if img.width != max_width:
            ratio = max_width / img.width
            new_height = int(img.height * ratio)
            img = img.resize((max_width, new_height), Image.Resampling.LANCZOS)
            print(f"Scaled image to full width: {max_width}x{new_height}")
        
        # Connect to printer with specified profile
        print(f"Connecting to printer at {printer_ip}:{printer_port} (Profile: {profile})")
        try:
            printer = Network(printer_ip, port=printer_port, profile=profile)
        except Exception as profile_error:
            print(f"Warning: Profile '{profile}' not found, using default profile")
            printer = Network(printer_ip, port=printer_port)
        
        # Print image only (no header/footer text)
        print("Printing image...")
        printer.image(img)
        printer.ln(1)  # Single line feed after image
        printer.cut()
        
        # Close connection
        printer.close()
        
        print("Image printed successfully!")
        return True
        
    except ConnectionError:
        print(f"Error: Could not connect to printer at {printer_ip}:{printer_port}")
        print("Make sure the printer is online and the IP address is correct")
        return False
    except Exception as e:
        print(f"Error printing image: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Print an image to a network printer using ESC/POS",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python print_image.py 192.168.1.100 image.jpg
  python print_image.py 192.168.1.100 image.png --port 9100
  python print_image.py 192.168.1.100 image.png --profile TM-T20II
  python print_image.py 192.168.1.100 image.png --width 576
        """
    )
    
    parser.add_argument("printer_ip", help="IP address of the printer")
    parser.add_argument("image_path", help="Path to the image file to print")
    parser.add_argument("--port", "-p", type=int, default=9100, 
                       help="Printer port (default: 9100)")
    parser.add_argument("--profile", type=str, default="TM-T20II",
                       help="Printer profile (default: TM-T20II for TM-T20III compatibility)")
    parser.add_argument("--width", type=int, default=576,
                       help="Maximum image width in pixels (default: 576 for 80mm paper)")
    
    args = parser.parse_args()
    
    # Print configuration
    print("ESC/POS Image Printer")
    print("=" * 20)
    print(f"Printer IP: {args.printer_ip}")
    print(f"Printer Port: {args.port}")
    print(f"Printer Profile: {args.profile}")
    print(f"Image: {args.image_path}")
    print()
    
    # Print the image
    success = print_image_to_network_printer(
        args.image_path, 
        args.printer_ip, 
        args.port,
        args.profile,
        args.width
    )
    
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main() 