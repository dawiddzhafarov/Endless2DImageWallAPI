package com.example.endless2dimagewall;

import androidx.appcompat.app.AppCompatActivity;
import androidx.constraintlayout.widget.ConstraintLayout;

import android.annotation.SuppressLint;
import android.os.Bundle;
import android.view.MotionEvent;
import android.view.View;
import android.view.Window;
import android.widget.FrameLayout;
import android.widget.LinearLayout;
import android.widget.RelativeLayout;

import java.lang.invoke.ConstantCallSite;

public class MainActivity extends AppCompatActivity{

    @SuppressLint("ClickableViewAccessibility")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        ConstraintLayout layout = (ConstraintLayout) findViewById(R.id.main);

        layout.setOnTouchListener((v, event) -> {
            float dx=0,dy=0,x=0,y=0;
            switch(event.getAction()) {
                case MotionEvent.ACTION_DOWN: {
                    x = event.getX();
                    y = event.getY();
                    dx = x - v.getX();
                    dy = y - v.getY();
                }
                break;
                case MotionEvent.ACTION_MOVE: {
                    v.setX(event.getX() - dx);
                    v.setY(event.getY() - dy);
                }
                break;
                case MotionEvent.ACTION_UP: {
                    //your stuff
                }
            }
            return true;
        });
    }



}