// +build !goci

// x c binding for the win
//
// connect to x server to:
//  * make a window
//  * get a cairo drawing surface
//  * get mouse and keyboard events

package eezl

/*
#cgo LDFLAGS: -lcairo -lxcb
#include <xcb/xcb.h>
#include <cairo/cairo-xcb.h>
#include <stdlib.h>
#include <string.h>

xcb_visualtype_t *get_root_visual_type(xcb_screen_t *s)
{
    xcb_visualtype_t *visual_type = NULL;
    xcb_depth_iterator_t depth_iter;

    depth_iter = xcb_screen_allowed_depths_iterator(s);

    for(;depth_iter.rem;xcb_depth_next(&depth_iter)) {
        xcb_visualtype_iterator_t visual_iter;

        visual_iter = xcb_depth_visuals_iterator(depth_iter.data);
        for(;visual_iter.rem;xcb_visualtype_next(&visual_iter)) {
            if(s->root_visual == visual_iter.data->visual_id) {
                visual_type = visual_iter.data;
                break;
            }
        }
    }

    return visual_type;
}
*/
import "C"

import (
	//"fmt"
	//"github.com/philetus/go-cairo"
)

// go struct to hold xcb connection data to an xserver screen
type Xscreen struct {
	conn *C.xcb_connection_t
	setup *C.xcb_setup_t
	screen *C.xcb_screen_t
	vistype *C.xcb_visualtype_t
}

// opens a connection to x server default screen
func Xconnect() *Xscreen {

	con := C.xcb_connect(nil, nil) // open connection to default x display
	set := C.xcb_get_setup(con)
	scr := C.xcb_setup_roots_iterator(set).data
	vzt := C.get_root_visual_type(scr)

	return &Xscreen{conn: con, setup: set, screen: scr, vistype: vzt}
}

type Xwin struct {
	xscreen *Xscreen
	height, width int
	//context_id C.xcb_gcontext_t
	window_id C.xcb_window_t
	//pixmap_id C.xcb_pixmap_t
}

func (self *Xscreen) NewXwin(hght, wdth int) *Xwin {
	//var pid C.xcb_pixmap_t = C.xcb_generate_id(xscr.conn)
	
	// generate context id and create graphics context 
	/*
	var cid C.xcb_gcontext_t = C.xcb_generate_id(self.conn)
	var cmsk int = C.XCB_GC_FOREGROUND | C.XCB_GC_BACKGROUND
	cval := []int{self.screen.black_pixel, self.screen.white_pixel}
	cvalp := (*C.int)(&cval[0])
	C.xcb_create_gc(self.conn, 
                    cid, 
					self.screen.root, 
					C.int(cmsk), 
					cvalp)
	*/

	// generate window id and create window
	var wid C.xcb_window_t = C.xcb_window_t(C.xcb_generate_id(self.conn))
	var wmsk int = C.XCB_CW_BACK_PIXEL | C.XCB_CW_EVENT_MASK
	wval := []C.uint32_t{self.screen.white_pixel, C.XCB_EVENT_MASK_EXPOSURE}
	wvalp := (*C.uint32_t)(&wval[0])
	C.xcb_create_window(self.conn,                          // connection
                        C.XCB_COPY_FROM_PARENT,             // depth
                        wid,                                // window id
                        self.screen.root,                   // parent window
                        0, 0,                               // x, y
                        C.uint16_t(wdth), C.uint16_t(hght), // width, height
                        0,                                  // border_width
                        C.XCB_WINDOW_CLASS_INPUT_OUTPUT,    // class
                        self.screen.root_visual,            // visual
                        C.uint32_t(wmsk), wvalp);                // masks
	
	// show window on screen
	C.xcb_map_window(self.conn, wid)
	C.xcb_flush(self.conn)

	return &Xwin{xscreen: self, height: hght, width: wdth, window_id: wid}
}

/*
func (self *Xwin) NewCairoSurface() *cairo.Surface {
	var csrf C.cairo_surface_t = C.cairo_xcb_surface_create(
									self.xscreen.conn,
									self.window_id,
									&xcb_visualtype, 
									C.int(wdth), C.int(hght))
	return cairo.NewSurfaceFromC(s *C.cairo_surface_t, c *C.cairo_t)
}
*/
