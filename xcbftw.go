// +build !goci

// x c binding for the window
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

// functions to cast generic events to specific events as go doesnt want to
xcb_configure_notify_event_t *cast_configure_notify_event(xcb_generic_event_t *event)
{
	return (xcb_configure_notify_event_t *)event;
}
xcb_expose_event_t *cast_expose_event(xcb_generic_event_t *event)
{
	return (xcb_expose_event_t *)event;
}
xcb_button_press_event_t *cast_button_press_event(xcb_generic_event_t *event)
{
	return (xcb_button_press_event_t *)event;
}
xcb_button_release_event_t *cast_button_release_event(xcb_generic_event_t *event)
{
	return (xcb_button_release_event_t *)event;
}
xcb_motion_notify_event_t *cast_motion_notify_event(xcb_generic_event_t *event)
{
	return (xcb_motion_notify_event_t *)event;
}
xcb_key_press_event_t *cast_key_press_event(xcb_generic_event_t *event)
{
	return (xcb_key_press_event_t *)event;
}
xcb_key_release_event_t *cast_key_release_event(xcb_generic_event_t *event)
{
	return (xcb_key_release_event_t *)event;
}
*/
import "C"

import (
	"fmt"
	"github.com/philetus/eezl/keys"
)

// go struct to hold xcb connection data to an xserver screen
type Xscreen struct {
	conn *C.xcb_connection_t
	setup *C.xcb_setup_t
	screen *C.xcb_screen_t
	vistype *C.xcb_visualtype_t
	eezldeks map[C.xcb_window_t]*Eezl
}

// opens a connection to x server default screen
func Xconnect() *Xscreen {

	con := C.xcb_connect(nil, nil) // open connection to default x display
	set := C.xcb_get_setup(con)
	scr := C.xcb_setup_roots_iterator(set).data
	vzt := C.get_root_visual_type(scr)

	xscrn := &Xscreen{conn: con, setup: set, screen: scr, vistype: vzt,
					  eezldeks: make(map[C.xcb_window_t]*Eezl)}
	
	// start event loop to feed inputs and gels into pipes
	go xscrn.event_loop()
	
	return xscrn
}

// input flavrs
const (
	PointerMotion int = iota
	PointerPress
	PointerRelease
	KeyPress
	KeyRelease
)

// represents pointer motion and press events and keyboard events
type Input struct {
	Flavr int
	Timestamp int
	Y, X int
	Stroke *keys.Key
}

// request to generate a new gel
type gelStain struct {
	Resize bool
	Height, Width int
}

// eezl is window with cairo drawing surface and event handlers
// 
// * input events are passed thru the input pipe
// * gels (individual 2d frames with cairo vector drawing methods) are passed
//   thru the gel pipe
// *
type Eezl struct {
	xscreen *Xscreen
	dead bool
	height, width int
	window_id C.xcb_window_t
	pixmap_id C.xcb_pixmap_t
	xcontext_id C.xcb_gcontext_t
	surface *C.cairo_surface_t
	stainPipe chan *gelStain
	
	InputPipe chan *Input
	GelPipe chan *Gel
}

func (self *Xscreen) NewEezl(hght, wdth int) *Eezl {	

	// generate window id and create window
	var wid C.xcb_window_t = C.xcb_window_t(C.xcb_generate_id(self.conn))
	var wmsk int = C.XCB_CW_BACK_PIXEL | C.XCB_CW_EVENT_MASK
	wval := []C.uint32_t{
				self.screen.white_pixel, 
					C.XCB_EVENT_MASK_EXPOSURE |
					C.XCB_EVENT_MASK_STRUCTURE_NOTIFY |
					C.XCB_EVENT_MASK_BUTTON_PRESS |
					C.XCB_EVENT_MASK_BUTTON_RELEASE |
					C.XCB_EVENT_MASK_POINTER_MOTION |
					C.XCB_EVENT_MASK_KEY_PRESS |
					C.XCB_EVENT_MASK_KEY_RELEASE}
					
	wvalp := (*C.uint32_t)(&wval[0])
	C.xcb_create_window(self.conn,                          // x connection
                        self.screen.root_depth,             // depth (must match pixmap)
                        wid,                                // window id
                        self.screen.root,                   // parent window
                        0, 0,                               // x, y
                        C.uint16_t(wdth), C.uint16_t(hght), // width, height
                        0,                                  // border_width
                        C.XCB_WINDOW_CLASS_INPUT_OUTPUT,    // class
                        self.screen.root_visual,            // visual
                        C.uint32_t(wmsk), wvalp);           // masks
	
	// generate pixmap for double-buffered rendering
	var pid C.xcb_pixmap_t = C.xcb_pixmap_t(C.xcb_generate_id(self.conn))
	C.xcb_create_pixmap(self.conn,                          // x connection
                        self.screen.root_depth,             // depth of the screen
                        pid,                                // id of the pixmap
                        C.xcb_drawable_t(self.screen.root), // ???
                        C.uint16_t(wdth), C.uint16_t(hght)) // width, height
                        	
	// create simple graphics context for copying pixmap buffer to window
	var xcid C.xcb_gcontext_t = C.xcb_gcontext_t(C.xcb_generate_id(self.conn))
	var xcmsk int = C.XCB_GC_FOREGROUND | C.XCB_GC_BACKGROUND
	xcval := []C.uint32_t{self.screen.black_pixel, self.screen.white_pixel}
	xcvalp := (*C.uint32_t)(&xcval[0])
	C.xcb_create_gc(self.conn, 
                    xcid, 
					C.xcb_drawable_t(self.screen.root), 
					C.uint32_t(xcmsk), 
					xcvalp)

	// create a cairo surface tied to pixmap for rendering to buffer
	srf := C.cairo_xcb_surface_create(self.conn,                // x connection
									  C.xcb_drawable_t(pid),    // drawable
									  self.vistype,             // visual type
									  C.int(wdth), C.int(hght)) // width, height
	
	// show eezl window on screen
	C.xcb_map_window(self.conn, wid)
	C.xcb_flush(self.conn)

	ezl := &Eezl{xscreen: self,
				 dead: false,
				 height: hght, width: wdth, 
				 window_id: wid, pixmap_id: pid, xcontext_id: xcid, 
				 surface: srf, 
				 stainPipe: make(chan *gelStain, 1), // stainpipe holds 1 request
				 InputPipe: make(chan *Input, 256),
				 GelPipe: make(chan *Gel)}
	
	// add new eezl to xscreen eezls map
	self.eezldeks[wid] = ezl
	
	// start loop to redraw window
	go ezl.stain_loop()
	
	return ezl
}

// process gel stains as they come off of stain pipe
func (self *Eezl) stain_loop() {
	for !self.dead {
		stn := <-self.stainPipe
		if stn.Resize {
		
			// free old surface and pixmap
			C.cairo_surface_finish(self.surface)
			C.xcb_free_pixmap(self.xscreen.conn, self.pixmap_id)
			
			// allocate new pixmap and surface with new size
			//self.pixmap_id = C.xcb_pixmap_t(C.xcb_generate_id(self.xscreen.conn)) ???
			C.xcb_create_pixmap(self.xscreen.conn,
                        		self.xscreen.screen.root_depth,
                        		self.pixmap_id,
                        		C.xcb_drawable_t(self.xscreen.screen.root),
                        		C.uint16_t(stn.Width), C.uint16_t(stn.Height))
			self.surface = C.cairo_xcb_surface_create(
							   self.xscreen.conn,
							   C.xcb_drawable_t(self.pixmap_id),
							   self.xscreen.vistype,
							   C.int(stn.Width), C.int(stn.Height))
		}
		
		// create cairo drawing context from cairo surface and fill with white
		cntxt := C.cairo_create(self.surface)
		C.cairo_set_source_rgba(cntxt, 1.0, 1.0, 1.0, 1.0)
		C.cairo_paint(cntxt)
		
		// get new gel and send it down gelpipe to be drawn to
		gel := &Gel{context: cntxt,
					trigger_pipe: make(chan bool, 1),
					Height: stn.Height, Width: stn.Width}
		self.GelPipe <- gel
		
		// block until trigger passed
		if <-gel.trigger_pipe {
		
			// destroy gels cairo context
			C.cairo_destroy(gel.context)
						
			// copy pixmap buffer from gel onto window
			C.xcb_copy_area(self.xscreen.conn,
							C.xcb_drawable_t(self.pixmap_id),
							C.xcb_drawable_t(self.window_id),
							self.xcontext_id,
							0, 0, 0, 0,
							C.uint16_t(gel.Width), C.uint16_t(gel.Height))
							
			C.xcb_flush(self.xscreen.conn)
		}
	}
}

func (self *Xscreen) event_loop() {
	for {
		evnt := C.xcb_wait_for_event(self.conn)
		switch evnt.response_type &^ 0x80 {
		
			case C.XCB_CONFIGURE_NOTIFY:
				cne := C.cast_configure_notify_event(evnt)
				wid := cne.window
				ezl := self.eezldeks[wid]
				h := int(cne.height)
				w := int(cne.width)
				
				// if window size has changed send resize stain 
				// and block until handled
				if h != ezl.height || w != ezl.width {
					ezl.height = h
					ezl.width = w
					ezl.stainPipe <- &gelStain{Resize: true, 
											   Height: h, Width: w}
				}

				//fmt.Printf("caught configure notify event for window %d!\n", int(wid))
				
			case C.XCB_EXPOSE:
				ee := C.cast_expose_event(evnt)
				wid := ee.window
				ezl := self.eezldeks[wid]
				
				// wrap sending dirty signal to eezl in select with default
				// to make it non-blocking -- if there is already a new gel 
				// pending ignore expose event
				select {
					case ezl.stainPipe <- &gelStain{Resize: false, 
													Height: ezl.height, 
													Width: ezl.width}:
					default:
				}
				//fmt.Printf("caught expose event for window %d!\n", int(wid))

			case C.XCB_BUTTON_PRESS:
				bpe := C.cast_button_press_event(evnt)
				wid := bpe.event
				inp := &Input{Flavr: PointerPress,
							  Timestamp: int(bpe.time),
							  Y: int(bpe.event_y),
							  X: int(bpe.event_x)}
				self.eezldeks[wid].InputPipe <- inp

			case C.XCB_BUTTON_RELEASE:
				bre := C.cast_button_release_event(evnt)
				wid := bre.event
				inp := &Input{Flavr: PointerRelease,
							  Timestamp: int(bre.time),
							  Y: int(bre.event_y),
							  X: int(bre.event_x)}
				self.eezldeks[wid].InputPipe <- inp

			case C.XCB_MOTION_NOTIFY:
				mne := C.cast_motion_notify_event(evnt)
				wid := mne.event
				inp := &Input{Flavr: PointerMotion,
							  Timestamp: int(mne.time),
							  Y: int(mne.event_y),
							  X: int(mne.event_x)}
				self.eezldeks[wid].InputPipe <- inp

			case C.XCB_KEY_PRESS:
				kpe := C.cast_key_press_event(evnt)
				wid := kpe.event
				inp := &Input{Flavr: KeyPress,
							  Timestamp: int(kpe.time),
							  Stroke: keys.New(int(kpe.detail))}
				self.eezldeks[wid].InputPipe <- inp

			case C.XCB_KEY_RELEASE:
				kre := C.cast_key_release_event(evnt)
				wid := kre.event
				inp := &Input{Flavr: KeyRelease,
							  Timestamp: int(kre.time),
							  Stroke: keys.New(int(kre.detail))}
				self.eezldeks[wid].InputPipe <- inp
			
		    // random unhelpful events?
		    case C.XCB_NO_EXPOSURE:
		    case C.XCB_MAP_NOTIFY:
		    case C.XCB_REPARENT_NOTIFY:
		    
			default:
				fmt.Printf("caught unknown event: %d!\n", int(evnt.response_type))
		}
	}
}
